from .utils.getDataset import download_file, read_file

datafile_path = download_file(
	"https://raw.githubusercontent.com/rasbt/LLMs-from-scratch/main/ch02/01_main-chapter-code/",
	"the-verdict.txt",
)
raw_text = read_file(datafile_path)
print(raw_text[:500])